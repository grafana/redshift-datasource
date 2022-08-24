import { DataFrame, DataQuery, DataQueryRequest, DataQueryResponse, LoadingState } from '@grafana/data';
import { Observable, Subscription } from 'rxjs';

export interface RequestLoopOptions<TQuery extends DataQuery = DataQuery> {
  /**
   * If the response needs an additional request to execute, return it here
   */
  getNextQuery: (rsp: DataQueryResponse) => TQuery | undefined;

  /**
   * The datasource execute method
   */
  query: (req: DataQueryRequest<TQuery>) => Observable<DataQueryResponse>;

  /**
   * Process the results
   */
  process: (data: DataFrame[]) => DataFrame[];

  /**
   * Check if the query should be cancelled
   */
  shouldCancel: () => boolean;

  /**
   * Callback that gets executed when unsubscribed
   */
  onCancel: () => void;
}

/**
 * Continue executing requests as long as `getNextQuery` returns a query
 */
export function getRequestLooper<T extends DataQuery = DataQuery>(
  req: DataQueryRequest<T>,
  options: RequestLoopOptions<T>
): Observable<DataQueryResponse> {
  return new Observable<DataQueryResponse>((subscriber) => {
    let nextQuery: T | undefined = undefined;
    let subscription: Subscription | undefined = undefined;
    let loadingState: LoadingState | undefined = LoadingState.Loading;
    let nextRequestDelay = 1; // Seconds until the next request
    let count = 1;
    let shouldCancel = false;

    // Single observer gets reused for each request
    const observer = {
      next: (rsp: DataQueryResponse) => {
        loadingState = rsp.state;
        let checkstate = false;
        if (loadingState !== LoadingState.Error) {
          nextQuery = options.getNextQuery(rsp);
          const _shouldCancel = options.shouldCancel();

          if (nextQuery && _shouldCancel) {
            // `shouldCancel` is set here only if there is a `nextQuery`, otherwise
            // we would try to cancel a finished query in the cleanup function
            shouldCancel = _shouldCancel;
            nextQuery = undefined;
          }

          checkstate = true;
        }
        const data = options.process(rsp.data);

        // Show the spinner or streaming (streaming will show data)
        if (checkstate) {
          if (nextQuery) {
            if (data.length && data[0].length) {
              loadingState = LoadingState.Streaming;
            } else {
              loadingState = LoadingState.Loading;
            }
            // Calculate the time for the next request, cap at 10s
            nextRequestDelay = nextRequestDelay * 2 > 10 ? 10 : nextRequestDelay * 2;
          } else {
            loadingState = LoadingState.Done;
            nextRequestDelay = 0;
          }
        }
        subscriber.next({ ...rsp, data, state: loadingState, key: req.requestId });
      },
      error: (err: any) => {
        subscriber.error(err);
      },
      complete: () => {
        // Completion of one query
        if (subscription) {
          subscription.unsubscribe();
          subscription = undefined;
        }

        // Let the previous request finish first
        if (nextQuery) {
          const next = nextQuery;
          setTimeout(() => {
            subscription = options
              .query({ ...req, requestId: `${req.requestId}.${++count}`, targets: [next] })
              .subscribe(observer);
            nextQuery = undefined;
          }, nextRequestDelay * 1000);
        } else {
          subscriber.complete();
        }
      },
    };

    // First request
    subscription = options.query(req).subscribe(observer);

    return () => {
      observer.complete();
      if (nextQuery || shouldCancel) {
        options.onCancel();
      }
      nextQuery = undefined;
    };
  });
}
