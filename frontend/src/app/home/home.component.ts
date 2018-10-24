import { Component, OnInit } from '@angular/core';
import { AuthService } from '../auth/auth.service';
import { Corporation } from '../corporation/corporation';
import { CorporationService } from '../corporation/corporation.service';

@Component({
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  corporation: Corporation;
  corporationLogo: string;

  constructor(private auth: AuthService,
    private corporationService: CorporationService) {

  }

  ngOnInit() {
    if (this.auth.isLoggedIn()) {
      this.corporationService.corporation.subscribe(corporation => {
        this.corporation = corporation;
        this.corporationLogo = this.corporationService.getCorporationLogo(corporation.corporationID);
      });
    }
  }

}
