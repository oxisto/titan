import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { JwtModule } from '@auth0/angular-jwt';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { AppComponent } from './app.component';
import { AuthService } from './auth/auth.service';
import { BlueprintDetailComponent } from './blueprint/blueprint-detail.component';
import { BlueprintService } from './blueprint/blueprint.service';
import { BlueprintsComponent } from './blueprint/blueprints.component';
import { CharacterService } from './character/character.service';
import { CorporationService } from './corporation/corporation.service';
import { HomeComponent } from './home/home.component';
import { ManufacturingService } from './manufacturing/manufacturing.service';
import { MarketService } from './market/market.service';
import { ValuesPipe } from './values.pipe';

export function tokenGetter() {
  return localStorage.getItem('token');
}

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    BlueprintsComponent,
    BlueprintDetailComponent,
    ValuesPipe
  ],
  imports: [
    NgbModule,
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
