import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {Observable, Subscription} from 'rxjs';
import {BreakpointObserver} from '@angular/cdk/layout';
import {map, shareReplay} from 'rxjs/operators';
import {CommonService} from './service/commom.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})

export class AppComponent implements OnInit {

    public headerTitle = 'process manager';

    private subscription: Subscription;

    title = 'process manager';
    isHandset$: Observable<boolean> = this.breakpointObserver.observe('(max-width: 1200px)')
        .pipe(
            map(result => result.matches),
            shareReplay()
        );

    constructor(
        private breakpointObserver: BreakpointObserver,
        private commonService: CommonService,
        private cd: ChangeDetectorRef,
    ) {
    }

    ngOnInit(): void {
        this.subscription = this.commonService.sharedDataSource$.subscribe(
            msg => {
                console.log('[Sample1Component] shared data updated.');
                this.headerTitle = msg;
                this.cd.detectChanges();
            }
        );

    }
}
