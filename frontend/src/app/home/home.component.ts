import { Component, OnInit } from '@angular/core';
import { Corporation, Wallets } from '../corporation/corporation';
import { CorporationService } from '../corporation/corporation.service';

@Component({
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  corporation: Corporation;
  wallets: Wallets;
  corporationLogo: string;

  constructor(private corporationService: CorporationService) {

  }

  ngOnInit() {
    this.corporationService.getCorporation().subscribe(corporation => {
      this.corporation = corporation;
      this.corporationLogo = this.corporationService.getCorporationLogo(corporation.corporationID);
    });

    this.corporationService.getCorporationWallets().subscribe(wallets => {
      this.wallets = wallets;
    });
  }

}
