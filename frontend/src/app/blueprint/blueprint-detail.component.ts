import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ClipboardService } from 'ngx-clipboard';
import { BlueprintService } from './blueprint.service';

@Component({
  templateUrl: './blueprint-detail.component.html',
  styleUrls: ['./blueprint-detail.component.css'],
})
export class BlueprintDetailComponent implements OnInit {

  public manufacturing: any;

  private possibleME: number[] = Array.from(new Array(11), (x, i) => i);
  private possibleTE: number[] = Array.from(new Array(11), (x, i) => i * 2);

  ME = 0;
  TE = 0;

  typeID: number;

  constructor(
    private blueprintService: BlueprintService,
    private route: ActivatedRoute,
    private clipboardService: ClipboardService) {
  }

  ngOnInit() {
    this.route.params.subscribe(params => {
      this.typeID = +params['typeID'];

      this.updateType();
    });
  }

  updateType(): any {
    this.blueprintService.getManufacturing(this.typeID, this.ME, this.TE).subscribe((manufacturing: any) => {
      this.manufacturing = manufacturing;

      // lock ME and TE for tech2
      if (manufacturing.isTech2) {
        this.ME = manufacturing.me;
        this.TE = manufacturing.te;
      }
    });
  }

  onOptionsChanged(event) {
    console.log(event);

    this.updateType();
  }

  copyMaterialsToClipboard(materials: any) {
    let text = '';
    for (const material of Object.values(materials) as any[]) {
      text += material.quantity + ' ' + material.name.en + '\n';
    }

    this.clipboardService.copyFromContent(text);
  }

}
