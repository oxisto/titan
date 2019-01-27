import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { IndustryJob, IndustryJobs } from './industry';

@Injectable({
  providedIn: 'root'
})
export class IndustryService {

  constructor(private http: HttpClient) { }

  getJobs(): Observable<IndustryJobs> {
    return this.http.get<IndustryJobs>('/api/industry/jobs').pipe(map(data => {
      return Object.assign(new IndustryJobs(), {
        ...data,
        jobs: new Map(Object.entries(data.jobs).map(entry => {
          entry[1] = Object.assign(new IndustryJob(), {
            ...entry[1],
            endDate: new Date(entry[1].endDate * 1000),
            startDate: new Date(entry[1].startDate * 1000),
            completedDate: new Date(entry[1].completedDate * 1000),
            pauseDate: new Date(entry[1].pauseDate * 1000)
          });
          return entry;
        }))
      });
    }));
  }
}
