import { Component, Input, OnInit } from '@angular/core';
import { MarketService } from '../market/market.service';

@Component({
  selector: 'app-type',
  templateUrl: './type.component.html',
  styleUrls: ['./type.component.css']
})
export class TypeComponent implements OnInit {

  @Input() type: any;
  @Input() skill: boolean;
  @Input() blueprint: boolean;

  constructor(private marketService: MarketService) { }

  ngOnInit() {
  }

  openMarketView(typeID: number) {
    this.marketService.postOpenMarketView(typeID).subscribe(resp => { });
  }


}
