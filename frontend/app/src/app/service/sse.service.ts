import {Injectable, NgZone} from '@angular/core';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class SseService {

    private eventSource!: EventSource;

    constructor(private _zone: NgZone) {
    }

    private static getEventSource(url: string): EventSource {
        return new EventSource(url);
    }

    public closeServerSentEvent() {
        this.eventSource.close();
    }

    public getServerSentEvent(url: string): Observable<any> {
        return new Observable(observer => {
            this.eventSource = SseService.getEventSource(url);
            this.eventSource.onmessage = event => {
                this._zone.run(() => {
                    observer.next(event);
                });
            };
            this.eventSource.onerror = error => {
                this._zone.run(() => {
                    observer.error(error);
                    this.closeServerSentEvent()
                });
            };
        });
    }

}
