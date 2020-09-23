import {Component, OnInit} from '@angular/core';

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
        this.customHost = '';
    }

    public setHost(): void {
        console.log('location.hostname');
    }

}
