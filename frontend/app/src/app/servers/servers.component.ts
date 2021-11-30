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
    public localhostName = location.hostname;
    public inputIPAdder = '192.168.10.109';
    public inputPort = 5984;

    private serverList = [];
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
        this.http.get(`${config.httpScheme}${location.hostname}:${config.port}/calculator`).subscribe(
            (data: any) => {
                console.log(data)
                if (data == null) {
                    this.serverList = []
                } else {
                    this.serverList = data
                }
                for (let server of this.serverList) {
                    this.sseService.getServerSentEvent(
                        `${config.httpScheme}${server['IP']}:${server['Port']}/health`
                    ).subscribe((hostData: any) => {
                        const data: hostStatus = JSON.parse(hostData.data);
                        this.serverStatuses[server['IP']] = [
                            [data.VRAM, 1 - data.VRAM],
                            [data.RAM, 1 - data.RAM],
                            server['Port'],
                            server['Status'],
                        ];
                    }, error => {
                        console.log(error)
                    })
                }
            })
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
            const formData = new FormData();
            formData.append('mode', 'join');
            formData.append('ip', String(this.inputIPAdder));
            formData.append('port', String(this.inputPort));
            this.http.post(`${config.httpScheme}${location.hostname}:${config.port}/calculator`, formData).subscribe(
                () => {
                    window.location.reload();
                }, error => {
                    if (String(error.error).indexOf('1062') != -1) {
                        alert('既に登録されています。')
                    } else {
                        alert(error.error);
                    }
                });
        } else {
            //ipアドレス以外
            alert('ipアドレスではありません');
        }
    }

    public onStopServer(ip: string, port: string): void {
        const formData = new FormData();
        formData.append('mode', 'stop');
        formData.append('ip', ip);
        formData.append('port', port);
        this.http.post(`${config.httpScheme}${location.hostname}:${config.port}/calculator`, formData).subscribe(
            (data: any) => {
                console.log(data);
                window.location.reload();
            }
        )
    }

    public onResumeServer(ip: string, port: string): void {
        const formData = new FormData();
        formData.append('mode', 'active');
        formData.append('ip', ip);
        formData.append('port', port);
        this.http.post(`${config.httpScheme}${location.hostname}:${config.port}/calculator`, formData).subscribe(
            (data: any) => {
                console.log(data);
                window.location.reload();
            }, error => {
                alert(`接続に失敗しました: ${error.error}`)
            }
        )
    }

    public onDeleteServer(ip: string): void {
        this.http.delete(`${config.httpScheme}${location.hostname}:${config.port}/calculator?ip=${ip}`).subscribe(
            (data: any) => {
                console.log(data);
                window.location.reload();
            }
        )
    }

}
