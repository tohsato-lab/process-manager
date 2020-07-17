import {Component} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import config from '../../config';


@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})

export class AppComponent {

    public hiddenUploadPage = true;

    constructor(private http: HttpClient) {
    }

    public onSelect(event): void {
        console.log(event);
        for (const file of [...event.addedFiles]) {
            this.upload(file);
        }
    }

    private upload(file): void {
        console.log('upload');
        const formData = new FormData();
        formData.append('file', file, file.name);
        this.http.post(
            `${config.urlScheme}${config.host}:${config.port}/upload`, formData
        ).subscribe(value => {
            console.log(value);
        });
    }

    public onAddButton(): void {
        this.hiddenUploadPage = false;
    }

}
