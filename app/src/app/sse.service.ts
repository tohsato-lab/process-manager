import {Injectable, NgZone} from '@angular/core';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class SseService {

    constructor(private _zone: NgZone) {
    }

    private static getEventSource(url: string): EventSource {
        return new EventSource(url);
    }

    public getServerSentEvent(url: string): Observable<any> {
        return new Observable(observer => {
            const eventSource = SseService.getEventSource(url);
            eventSource.onmessage = event => {
                this._zone.run(() => {
                    observer.next(event);
                });
            };
            eventSource.onerror = error => {
                this._zone.run(() => {
                    observer.error(error);
                });
            };
        });
    }

}
