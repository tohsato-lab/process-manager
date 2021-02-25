import {Component, OnDestroy, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import config from '../../../config';
import {Subscription} from 'rxjs';
import {CommonService} from '../service/commom.service';
import {SseService} from '../service/sse.service';


@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit, OnDestroy {

    public hiddenUploadPage = true;
    public fileInfos: any = [];
    public processList = [];
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
        this.sseService.getServerSentEvent(
            `${config.httpScheme}${location.hostname}:${config.port}/process_status`
        ).subscribe((processData: any) => {
            this.processList = JSON.parse(processData.data);
            for (const process of this.processList) {
                process.Selected = false;
            }
            console.log(this.processList);
        });
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

    public onOpenExplorer(id) {
        window.location.href = `${config.httpScheme}${location.hostname}:${config.port}/programs/${id}`;
    }

    public onAddButton(): void {
        this.hiddenUploadPage = false;
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/env_info`
        ).subscribe(value => {
            this.envList = value;
        }, error => {
            console.log(error);
        });
    }

    public onCloseUpload(): void {
        this.hiddenUploadPage = true;
        this.fileInfos = [];
    }

    public onSelectFiles(event): void {
        console.log(event);
        for (const file of [...event.addedFiles]) {
            this.fileInfos.push({file: file, vram: 0.0, env: 'base', target: 'main.py', exec_count: 1});
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

    public onDelete(id): void {
        this.http.get(
            `${config.httpScheme}${location.hostname}:${config.port}/delete?id=${id}`
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }

    public onUpload(): void {
        for (const fileInfo of this.fileInfos) {
            this.upload(fileInfo);
        }
        this.fileInfos = [];
    }

    private upload(info): void {
        console.log('upload');
        const formData = new FormData();
        formData.append('file', info.file, info.file.name);
        formData.append('vram', info.vram);
        formData.append('env', info.env);
        formData.append('target', info.target);
        formData.append('exec_count', info.exec_count);
        formData.append('comment', 'テスト');
        this.onCloseUpload();
        this.http.post(
            `${config.httpScheme}${location.hostname}:${config.port}/upload`, formData
        ).subscribe(value => {
            console.log(value);
        }, error => {
            alert(error.error.text);
        });
    }
}
