import { Component, OnInit } from '@angular/core';
import { Headers, Http } from '@angular/http';

import 'rxjs/add/operator/map';

import { AuthService } from '../auth/auth.service';
import { CorporationService } from '../corporation/corporation.service';
import { ManufacturingService } from '../manufacturing/manufacturing.service';
import { Corporation } from '../corporation/corporation';

@Component({
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  corporation: Corporation;

  constructor(private auth: AuthService,
              private corporationService: CorporationService,
              private manufacturingService: ManufacturingService) {

  }

  ngOnInit() {
    if (this.auth.isLoggedIn()) {
      this.corporationService.corporation.subscribe(corporation => this.corporation = corporation);
    }
  }

}
