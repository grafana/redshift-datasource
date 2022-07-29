import { applySQLTemplateVariables, filterSQLQuery } from '@grafana/aws-sdk';
import { DataFrame, DataQueryRequest, DataQueryResponse, DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import {
  DataSourceWithBackend,
  getTemplateSrv,
  getBackendSrv,
  toDataQueryResponse,
  BackendDataSourceResponse,
} from '@grafana/runtime';
import { getRequestLooper } from 'requestLooper';
import { merge, Observable, of } from 'rxjs';
import { map } from 'rxjs/operators';
import { RedshiftVariableSupport } from 'variables';

import { RedshiftCustomMeta, RedshiftDataSourceOptions, RedshiftQuery } from './types';

let requestCounter = 100;

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  private runningQueries: { [hash: string]: string };

  constructor(instanceSettings: DataSourceInstanceSettings<RedshiftDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new RedshiftVariableSupport(this);
    this.runningQueries = {};
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  filterQuery = filterSQLQuery;

  applyTemplateVariables = (query: RedshiftQuery, scopedVars: ScopedVars) =>
    applySQLTemplateVariables(query, scopedVars, getTemplateSrv);

  query(request: DataQueryRequest<RedshiftQuery>): Observable<DataQueryResponse> {
    const { intervalMs, maxDataPoints } = request;
    const targets = request.targets.filter(this.filterQuery).map((q) => ({
      ...q,
      intervalMs,
      maxDataPoints,
      datasource: this.getRef(),
      ...this.applyTemplateVariables(q, request.scopedVars),
    }));
    if (!targets.length) {
      return of({ data: [] });
    }
    const all: Array<Observable<DataQueryResponse>> = [];
    for (let target of targets) {
      if (target.hide) {
        continue;
      }
      all.push(this.doSingle(target, request));
    }
    return merge(...all);
  }

  storeQuery(target: RedshiftQuery, queryID: string) {
    const key = JSON.stringify(target);
    this.runningQueries[key] = queryID;
  }

  getQuery(target: RedshiftQuery) {
    const key = JSON.stringify(target);
    return this.runningQueries[key];
  }

  removeQuery(target: RedshiftQuery) {
    const key = JSON.stringify(target);
    delete this.runningQueries[key];
  }

  doSingle(target: RedshiftQuery, request: DataQueryRequest<RedshiftQuery>): Observable<DataQueryResponse> {
    let queryId: string | undefined = undefined;
    let allData: DataFrame[] = [];
    return getRequestLooper(
      { ...request, targets: [target], requestId: `aws_ts_${requestCounter++}` },
      {
        // Check for a "queryID" in the response
        getNextQuery: (rsp: DataQueryResponse) => {
          if (rsp.data?.length) {
            const first = rsp.data[0] as DataFrame;
            const meta = first.meta?.custom as RedshiftCustomMeta;
            if (meta && meta.queryID) {
              queryId = meta.queryID;
              this.storeQuery(target, meta.queryID);
              const status = meta.status;
              const notFinished = status === 'submitted' || status === 'running';
              return {
                ...target,
                queryID: meta.queryID,
                skipCache: notFinished,
              } as RedshiftQuery;
            }
          }
          this.removeQuery(target);
          return undefined;
        },

        /**
         * The original request
         */
        query: (request: DataQueryRequest<RedshiftQuery>) => {
          const { range, targets, requestId } = request;
          const [target] = targets;
          const { skipCache, ...query } = target;
          const data = {
            queries: [query],
            range: range,
            from: range.from.valueOf().toString(),
            to: range.to.valueOf().toString(),
          };

          let headers = {};
          if (skipCache) {
            headers = {
              'X-Cache-Skip': true,
            };
          }
          const options = {
            method: 'POST',
            url: '/api/ds/query',
            data,
            requestId,
            headers,
          };

          return getBackendSrv()
            .fetch<BackendDataSourceResponse>(options)
            .pipe(map((result) => result.data))
            .pipe(
              map((r) => {
                const frames = toDataQueryResponse({ data: r }).data as DataFrame[];
                return { data: frames };
              })
            );
        },

        /**
         * Process the results
         */
        process: (data: DataFrame[]) => {
          for (const frame of data) {
            if (frame.fields.length > 0) {
              allData.push(frame);
            }
          }

          return allData;
        },

        /**
         * Callback that gets executed when unsubscribed
         */
        onCancel: () => {
          if (queryId) {
            this.removeQuery(target);
            this.postResource('cancel', {
              queryId,
            })
              .then((v) => {
                console.log('Query canceled:', v);
              })
              .catch((err) => {
                err.isHandled = true; // avoid the popup
                console.log('error killing', err);
              });
          }
        },
      }
    );
  }

  async cancel(target: RedshiftQuery) {
    const queryId = this.getQuery(target);
    try {
      this.removeQuery(target);
      await this.postResource('cancel', { queryId });
    } catch (err: any) {
      err.isHandled = true; // avoid the popup
      console.log('error killing', err);
    }
  }
}
