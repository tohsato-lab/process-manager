import {Component, OnDestroy, OnInit} from '@angular/core';
import {SseService} from '../service/sse.service';
import config from '../../../config';
import {CommonService} from '../service/commom.service';
import {Subscription} from 'rxjs';

@Component({
    selector: 'app-hosts',
    templateUrl: './hosts.component.html',
    styleUrls: ['./hosts.component.css']
})
export class HostsComponent implements OnInit, OnDestroy {

    public ramStatusList = {};
    public IPList = [location.hostname];

    private headerTitle = 'ホスト一覧';
    private subscription: Subscription;

    constructor(
        private sseService: SseService,
        private commonService: CommonService,
    ) {
    }

    ngOnInit(): void {
        this.commonService.onNotifySharedDataChanged(this.headerTitle);
        for (let ip of this.IPList) {
            this.sseService.getServerSentEvent(`${config.httpScheme}${ip}:${config.port}/gpu_status`)
                .subscribe(data => {
                    this.ramStatusList[ip] = data.data;
                });
        }
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        this.subscription.unsubscribe();
    }

}
