import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { BlueprintService } from './blueprint.service';

@Component({
  templateUrl: './blueprint-detail.component.html',
  styleUrls: ['./blueprint-detail.component.css'],
})
export class BlueprintDetailComponent implements OnInit {

  public manufacturing: any;

  private possibleME: number[] = Array.from(new Array(11), (x, i) => i);
  private possibleTE: number[] = Array.from(new Array(11), (x, i) => i * 2);

  constructor(private blueprintService: BlueprintService,
    private route: ActivatedRoute) {

  }

  ngOnInit() {
    this.route.params.subscribe(params => {
      const typeID = +params['typeID'];

      this.blueprintService.getManufacturing(typeID).subscribe(manufacturing => {
        this.manufacturing = manufacturing;
      });
    });
  }

}
