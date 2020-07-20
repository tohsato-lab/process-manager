import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';

@Component({
    selector: 'app-menu',
    templateUrl: './menu.component.html',
    styleUrls: ['./menu.component.css']
})
export class MenuComponent implements OnInit {

    constructor(private router: Router) {
    }

    ngOnInit(): void {
    }

    public onHome(): void {
        this.router.navigate(['home']);
    }

    public onDataset(): void {
        this.router.navigate(['dataset']);
    }

    public onHistory(): void {
        this.router.navigate(['history']);
    }

}
