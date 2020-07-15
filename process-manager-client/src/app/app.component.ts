import {Component} from '@angular/core';
import {HttpClient} from '@angular/common/http';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent {

  title = 'dropzone';
  files: File[] = [];

  constructor(private http: HttpClient) {
  }


  public onSelect(event): void {
    console.log(event);
    this.files.push(...event.addedFiles);
    const formData = new FormData();
    for (const file of this.files) {
      formData.append('file', file);
    }

    this.upload(formData);

  }

  public onRemove(event): void {
    console.log(event);
    this.files.splice(this.files.indexOf(event), 1);

  }

  private upload(formData): void {
    console.log('upload');
    /*
    if (files.length <= 0) {
      return;
    }

    const file = files[0];
    const formData = new FormData();

    formData.append('file', file, file.name);
     */

    this.http.post('http://localhost:8081/upload', formData)
      .subscribe(value => {
        console.log(value);
      });
  }

}
