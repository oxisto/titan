import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { BlueprintsComponent } from './blueprint/blueprints.component';
import { BlueprintDetailComponent } from './blueprint/blueprint-detail.component';
import { AuthService } from './auth/auth.service';
import { BlueprintService } from './blueprint/blueprint.service';
import { CorporationService } from './corporation/corporation.service';
import { CharacterService } from './character/character.service';
import { ManufacturingService } from './manufacturing/manufacturing.service';
import { SidebarComponent } from './sidebar/sidebar.component';
import { MarketService } from './market/market.service';

import { KeysPipe } from './keys.pipe';
import { ValuesPipe } from './values.pipe';
import { HttpClientModule } from '@angular/common/http';
import { JwtModule } from '@auth0/angular-jwt';

export function tokenGetter() {
  return localStorage.getItem('token');
}

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    SidebarComponent,
    BlueprintsComponent,
    BlueprintDetailComponent,
    KeysPipe,
    ValuesPipe
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    JwtModule.forRoot({
      config: {
        tokenGetter: tokenGetter,
      }
    }),
    FormsModule,
    RouterModule.forRoot([{
      path: '',
      component: HomeComponent
    }, {
      path: 'home',
      component: HomeComponent
    }, {
      path: 'manufacturing',
      component: BlueprintsComponent
    }, {
      path: 'manufacturing/:typeID',
      component: BlueprintDetailComponent
    }
    ], { useHash: true })
  ],
  providers: [
    AuthService,
    CorporationService,
    CharacterService,
    BlueprintService,
    ManufacturingService,
    MarketService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
