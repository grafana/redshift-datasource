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
const runningStatuses = ['started', 'submitted', 'running'];
const isRunning = (status = '') => runningStatuses.includes(status);

export class DataSource extends DataSourceWithBackend<RedshiftQuery, RedshiftDataSourceOptions> {
  private runningQueries: { [hash: string]: { queryID: string; status: string } };

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

  storeQuery(target: RedshiftQuery, queryID: string, status = '') {
    const key = JSON.stringify(target);
    this.runningQueries[key] = { queryID, status };
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
    let queryID: string | undefined = undefined;
    let allData: DataFrame[] = [];
    return getRequestLooper(
      { ...request, targets: [target], requestId: `aws_ts_${requestCounter++}` },
      {
        // Check for a "queryID" in the response
        getNextQuery: (rsp: DataQueryResponse) => {
          if (rsp.data?.length) {
            const first = rsp.data[0] as DataFrame;
            const meta = first.meta?.custom as RedshiftCustomMeta;
            if (meta && isRunning(meta.status)) {
              const status = meta.status;
              queryID = meta.queryID;
              this.storeQuery(target, queryID, status);
              return {
                ...target,
                queryID,
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
          const [query] = targets;
          const data = {
            queries: [query],
            range: range,
            from: range.from.valueOf().toString(),
            to: range.to.valueOf().toString(),
          };

          let headers = {};
          const queryInfo = this.getQuery(target);
          if (queryInfo && isRunning(queryInfo.status)) {
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
          if (queryID) {
            this.removeQuery(target);
            this.postResource('cancel', {
              queryID,
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
    const queryID = this.getQuery(target);
    try {
      this.removeQuery(target);
      await this.postResource('cancel', { queryID });
    } catch (err: any) {
      err.isHandled = true; // avoid the popup
      console.log('error killing', err);
    }
  }
}
