import {Component, OnInit} from '@angular/core';
import {SseService} from '../service/sse.service';
import config from '../../../config';

@Component({
    selector: 'app-hosts',
    templateUrl: './hosts.component.html',
    styleUrls: ['./hosts.component.css']
})
export class HostsComponent implements OnInit {

    constructor(
        private sseService: SseService,
    ) {
    }

    ngOnInit(): void {
        this.sseService.getServerSentEvent(`${config.httpScheme}${location.hostname}:${config.port}/gpu_status`)
            .subscribe(data => console.log(data.data));
    }

}
