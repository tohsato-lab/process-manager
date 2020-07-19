import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import config from '../../config';

import {webSocket} from 'rxjs/webSocket';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})

export class AppComponent implements OnInit {

    public hiddenUploadPage = true;
    public files: any = [];

    constructor(private http: HttpClient) {
    }

    ngOnInit(): void {
        webSocket('ws://localhost:8081/process_status').subscribe(
            (message: any) => {
                console.log(message);
            },
            err => console.log(err),
            () => console.log('complete')
        );
    }

    public onAddButton(): void {
        this.hiddenUploadPage = false;
    }

    public onCloseUpload(): void {
        this.hiddenUploadPage = true;
        this.files = [];
    }

    public onSelect(event): void {
        console.log(event);
        for (const file of [...event.addedFiles]) {
            this.files.push({data: file, vram: 0.0});
        }
    }

    public onUpload(): void {
        for (const file of this.files) {
            this.upload(file.data, file.vram);
        }
        this.files = [];
    }

    private upload(file, vram): void {
        console.log('upload');
        const formData = new FormData();
        formData.append('file', file, file.name);
        formData.append('vram', vram);
        this.http.post(
            `${config.urlScheme}${config.host}:${config.port}/upload`, formData
        ).subscribe(value => {
            console.log(value);
        }, error => {
            console.log(error);
        });
    }
}
