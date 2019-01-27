import { Component, OnInit } from '@angular/core';
import { IndustryService } from './industry.service';

@Component({
  selector: 'app-industry-jobs',
  templateUrl: './industry-jobs.component.html',
  styleUrls: ['./industry-jobs.component.css']
})
export class IndustryJobsComponent implements OnInit {

  jobs: any;

  activityToName = {
    '1': 'Manufacturing',
    '2': ' Researching Technology',
    '3': ' Researching Time Productivity',
    '4': ' Researching Material Productivity',
    '5': ' Copying',
    '6': ' Duplicating',
    '7': ' Reverse Engineering',
    '8': ' Invention',
  };

  constructor(private industryService: IndustryService) { }

  ngOnInit() {
    this.industryService.getJobs().subscribe(jobs => {
      this.jobs = jobs;
    });
  }

}
