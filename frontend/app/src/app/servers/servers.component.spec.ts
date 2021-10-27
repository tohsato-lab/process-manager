import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ServersComponent } from './servers.component';

describe('HostsComponent', () => {
  let component: ServersComponent;
  let fixture: ComponentFixture<ServersComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ ServersComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ServersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
