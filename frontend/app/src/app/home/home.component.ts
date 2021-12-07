import {Component, OnDestroy, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import config from '../../../config';
import {Subscription} from 'rxjs';
import {CommonService} from '../service/commom.service';
import {SseService} from '../service/sse.service';
import {webSocket} from 'rxjs/webSocket';


@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit, OnDestroy {

    public hiddenUploadPage = true;
    public uploadInfos: any = [];
    public processList = [];
    public execEnvs = {};
    public envList: any;

    private subscription: Subscription;
    private headerTitle = 'プロセス一覧';

    constructor(
        private http: HttpClient,
        private commonService: CommonService,
        private sseService: SseService,
    ) {
    }

    ngOnInit(): void {
        webSocket(`${config.websocketScheme}${location.hostname}:${config.port}/process_status`).subscribe(
            (message: any) => {
                console.log(message);
                this.processList = message !== null ? message : [];
                for (const process of this.processList) {
                    process.Selected = false;
                }
                console.log(this.processList);
            },
            err => {
                console.log(err);
                console.log(`${config.httpScheme}${location.hostname}:${config.port}`);
            },
            () => console.log('complete')
        );
        this.commonService.onNotifySharedDataChanged(this.headerTitle);
        window.onbeforeunload = () => this.ngOnDestroy();
    }

    ngOnDestroy(): void {
        //  リソースリーク防止のため CommonService から subcribe したオブジェクトを破棄する
        if (this.subscription) {
            this.subscription.unsubscribe();
        }
        this.sseService.closeServerSentEvent();
    }

    getKeys(dict: any) {
        return Object.keys(dict)
    }

    public onOpenExplorer(id) {
        window.location.href = `${config.httpScheme}${location.hostname}:${config.port}/programs/${id}`;
    }

    public onAddButton(): void {
        this.hiddenUploadPage = false;
        this.http.get(`${config.httpScheme}${location.hostname}:${config.port}/conda`).subscribe(
            (data: any) => {
                console.log(data);
                this.execEnvs = data;
            }, error => {
                console.log(error);
                this.onCloseUpload();
                alert(error.error)
            }
        )
    }

    public onCloseUpload(): void {
        this.hiddenUploadPage = true;
        this.uploadInfos = [];
    }

    public onSelectFiles(event): void {
        console.log(event);
        for (const file of [...event.addedFiles]) {
            this.uploadInfos.push({
                file: file,
                vram: 0.0,
                env: this.execEnvs[this.getKeys(this.execEnvs)[0]]['Envs'][0],
                target: 'main.py',
                exec_count: 1,
                ip: this.getKeys(this.execEnvs)[0],
                comment: '',
            });
        }
    }

    public onKill(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/kill?id=${id}`
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }

    public onTrash(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/trash?id=${id}`
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }

    public onUpload(): void {
        for (const fileInfo of this.uploadInfos) {
            this.upload(fileInfo);
        }
        this.onCloseUpload();
    }

    private upload(info): void {
        console.log('upload');
        const formData = new FormData();
        formData.append('file', info.file, info.file.name);
        formData.append('conda_env', info.env);
        formData.append('target_file', info.target);
        formData.append('exec_count', info.exec_count);
        this.http.put(`${config.httpScheme}${info.ip}:${this.execEnvs[info.ip]['Port']}/upload`, formData
        ).subscribe((processIDs: Array<string>) => {
            const formData = new FormData();
            formData.append('process_ids', JSON.stringify(processIDs));
            formData.append('process_name', String(info.file.name).split('.')[0]);
            formData.append('conda_env', info.env);
            formData.append('server_ip', info.ip);
            formData.append('comment', info.comment);
            this.http.put(`${config.httpScheme}${location.hostname}:${config.port}/process`, formData
            ).subscribe(value => {
                console.log(value);
            },error => {
                alert(error.error);
            })
        }, error => {
            alert(error.error);
        });
    }
}
