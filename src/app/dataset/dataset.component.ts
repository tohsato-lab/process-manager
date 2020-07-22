import {Component, OnInit} from '@angular/core';
import config from '../../../config';

@Component({
    selector: 'app-dataset',
    templateUrl: './dataset.component.html',
    styleUrls: ['./dataset.component.css']
})
export class DatasetComponent implements OnInit {

    public customHost = '';

    constructor() {
    }

    ngOnInit(): void {
        this.customHost = config.host;
    }

    public setHost(): void {
        config.host = this.customHost;
    }

}
