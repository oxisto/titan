import { Component, OnInit } from '@angular/core';
import { Headers, Http } from '@angular/http';
import { ActivatedRoute } from '@angular/router';

import 'rxjs/add/operator/map';

import { AuthService } from '../auth/auth.service';
import { BlueprintService } from './blueprint.service';
import { Blueprint } from './blueprint';
import { MarketService } from '../market/market.service';

@Component({
  templateUrl: './blueprint-detail.component.html',
  styleUrls: ['./blueprint-detail.component.css'],
})
export class BlueprintDetailComponent implements OnInit {

  public manufacturing: any;

  private possibleME: number[] = Array.from(new Array(11), (x , i) => i);
  private possibleTE: number[] = Array.from(new Array(11), (x , i) => i * 2);

  constructor(private auth: AuthService,
              private blueprintService: BlueprintService,
              private marketService: MarketService,
              private route: ActivatedRoute) {

  }

  openMarketView(typeID: number) {
    this.marketService.postOpenMarketView(typeID).subscribe(resp => {});
  }

  ngOnInit() {
    if (this.auth.isLoggedIn()) {
      this.route.params.subscribe(params => {
        const typeID = +params['typeID'];

        this.blueprintService.getManufacturing(typeID).subscribe(manufacturing => {
          this.manufacturing = manufacturing;
        });
      });
    }
  }

}
