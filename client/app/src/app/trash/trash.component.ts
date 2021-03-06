import {Component, OnDestroy, OnInit} from '@angular/core';
import {CommonService} from '../service/commom.service';
import {Subscription} from 'rxjs';
import {SseService} from '../service/sse.service';
import config from '../../../config';
import {HttpClient} from '@angular/common/http';

@Component({
    selector: 'app-trash',
    templateUrl: './trash.component.html',
    styleUrls: ['./trash.component.css']
})
export class TrashComponent implements OnInit, OnDestroy {

    public processList = [];

    private headerTitle = 'ゴミ箱';
    private subscription: Subscription;

    constructor(
        private commonService: CommonService,
        private sseService: SseService,
        private http: HttpClient,
    ) {
    }

    private getTrashProcesses() {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/trash`
        ).subscribe((value: any) => {
            console.log(value)
            this.processList = value == null ? [] : value;
        }, error => {
            console.log(error);
        });
    }

    ngOnInit(): void {
        this.commonService.onNotifySharedDataChanged(this.headerTitle);
        this.getTrashProcesses();
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        if (this.subscription) {
            this.subscription.unsubscribe();
        }
        this.sseService.closeServerSentEvent();
    }

    public onOpenExplorer(id) {
        window.location.href = `${config.httpScheme}${location.hostname}:${config.port}/programs/${id}`;
    }

    public onRecover(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/trash?id=${id}`
        ).subscribe(value => {
            console.log(value);
            this.getTrashProcesses();
        }, error => {
            console.log(error);
        });
    }

    public onDelete(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/delete?id=${id}`
        ).subscribe(value => {
            console.log(value);
            this.getTrashProcesses();
        }, error => {
            console.log(error);
        });
    }

}
