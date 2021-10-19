import {Component, OnDestroy, OnInit} from '@angular/core';
import {SseService} from '../service/sse.service';
import config from '../../../config';
import {CommonService} from '../service/commom.service';
import {Subscription} from 'rxjs';
import {MultiDataSet, Label, Colors} from 'ng2-charts';
import {ChartType} from 'chart.js';
import {HttpClient} from '@angular/common/http';

interface hostStatus {
    RAM: number;
    VRAM: number;
}

@Component({
    selector: 'app-hosts',
    templateUrl: './servers.component.html',
    styleUrls: ['./servers.component.css']
})


export class ServersComponent implements OnInit, OnDestroy {

    public chartLabels: Label[] = ['used', 'free'];
    public chartColors: Colors[] = [{
        backgroundColor: ['#89c148', '#8d8d8d']
    }];
    public options = {
        rotation: Math.PI,
        circumference: Math.PI,
        tooltips: {enabled: false},
        hover: {mode: null},
    };
    public chartType: ChartType = 'doughnut';
    public serverStatuses: { [ip: string]: MultiDataSet } = {};
    public hiddenRegisterServer = true;
    public inputIPAdder = '';
    public inputPort = -1;

    private IPList: string[] = [location.hostname];
    private headerTitle = 'サーバーリスト';
    private subscription: Subscription;

    constructor(
        private sseService: SseService,
        private commonService: CommonService,
        private http: HttpClient,
    ) {
    }

    getKeys(data): any {
        return Object.keys(data);
    }

    ngOnInit(): void {
        this.commonService.onNotifySharedDataChanged(this.headerTitle);
        for (let ip of this.IPList) {
            this.sseService.getServerSentEvent(`${config.httpScheme}${ip}:${config.port}/host_status`).subscribe((hostData: any) => {
                // this.hostStatuses[ip] = JSON.parse(hostData.data);
                const data: hostStatus = JSON.parse(hostData.data);
                this.serverStatuses[ip] = [
                    [data.VRAM, 1 - data.VRAM],
                    [data.RAM, 1 - data.RAM],
                ];
            });
        }
        window.onbeforeunload = () => this.ngOnDestroy();
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        if (this.subscription) {
            this.subscription.unsubscribe();
        }
        this.sseService.closeServerSentEvent();
    }

    public onRegisterCtrl(): void {
        this.hiddenRegisterServer = !this.hiddenRegisterServer;
    }

    public onRegisterServer(): void {
        if (this.inputIPAdder.match(/^\d{1,3}(\.\d{1,3}){3}$/) && this.inputPort != -1) {
            console.log(`${config.httpScheme}${this.inputIPAdder}:${this.inputPort}`);
            this.http.get(
                `${config.httpScheme}${this.inputIPAdder}:${this.inputPort}`
            ).subscribe(value => {
                console.log(value);
            }, error => {
                console.log(error);
            });
            //ipアドレス
            console.log('ipアドレスです');
        } else {
            //ipアドレス以外
            console.log('ipアドレスではありません');
        }
    }

}
