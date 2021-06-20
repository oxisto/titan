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
            endDate: Date.parse(entry[1].endDate),
            startDate: Date.parse(entry[1].startDate),
            completedDate: entry[1].completedDate != null ? Date.parse(entry[1].completedDate) : null,
            pauseDate: entry[1].pauseDate != null ? Date.parse(entry[1].pauseDate) : null,
            blueprint: {
              typeID: entry[1].blueprintTypeID,
              typeName: entry[1].blueprintTypeName,
            }
          });
          return entry;
        }))
      });
    }));
  }
}
