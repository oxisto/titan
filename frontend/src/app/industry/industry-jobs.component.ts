import { Component, OnInit } from '@angular/core';
import { IndustryService } from './industry.service';

@Component({
  selector: 'app-industry-jobs',
  templateUrl: './industry-jobs.component.html',
  styleUrls: ['./industry-jobs.component.css']
})
export class IndustryJobsComponent implements OnInit {

  private jobs: any;

  constructor(private industryService: IndustryService) { }

  private activityToName = {
    1: 'Manufacturing',
    8: 'Invention'
  };

  ngOnInit() {
    this.industryService.getJobs().subscribe(jobs => {
      this.jobs = jobs;
    });
  }

}
