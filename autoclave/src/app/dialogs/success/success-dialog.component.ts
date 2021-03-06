import { Component, OnInit, Inject } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material";

@Component({
  selector: "success-dialog",
  templateUrl: "./success-dialog.component.html",
  styleUrls: ["./success-dialog.component.scss"]
})
export class SuccessDialogComponent implements OnInit {
  constructor(
    private ref: MatDialogRef<SuccessDialogComponent>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      msg: string;
      firstName: string;
    }
  ) {
    setTimeout(() => {
      this.ref.close();
    }, 2000);
  }

  ngOnInit() {}

  close = () => {
    this.ref.close();
  };
}
