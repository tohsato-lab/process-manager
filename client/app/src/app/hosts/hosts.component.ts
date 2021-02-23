import {Component, OnDestroy, OnInit} from '@angular/core';
import {SseService} from '../service/sse.service';
import config from '../../../config';
import {CommonService} from '../service/commom.service';
import {Subscription} from 'rxjs';

interface hostStatus {
    RAM: number;
    VRAM: number;
}

@Component({
    selector: 'app-hosts',
    templateUrl: './hosts.component.html',
    styleUrls: ['./hosts.component.css']
})


export class HostsComponent implements OnInit, OnDestroy {

    public hostStatuses: { [ip: string]: hostStatus } = {};

    private IPList: string[] = [location.hostname];
    private headerTitle = 'ホスト一覧';
    private subscription: Subscription;

    constructor(
        private sseService: SseService,
        private commonService: CommonService,
    ) {
    }

    getKeys(data): any {
        return Object.keys(data);
    }

    ngOnInit(): void {
        this.commonService.onNotifySharedDataChanged(this.headerTitle);
        for (let ip of this.IPList) {
            this.sseService.getServerSentEvent(`${config.httpScheme}${ip}:${config.port}/host_status`).subscribe(hostData => {
                this.hostStatuses[ip] = JSON.parse(hostData.data);
            });
        }
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        this.subscription.unsubscribe();
    }

}
