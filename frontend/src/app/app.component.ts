import { Component, OnInit } from '@angular/core';

import 'rxjs/add/operator/toPromise';

import { AuthService } from './auth/auth.service';
import { CorporationService } from './corporation/corporation.service';
import { CharacterService } from './character/character.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  loggedIn = false;

  title = 'app works!';

  constructor(private authService: AuthService) {

  }

  ngOnInit() {

  }

}
