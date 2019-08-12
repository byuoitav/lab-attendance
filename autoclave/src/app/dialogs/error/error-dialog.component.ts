import { Component, OnInit, Inject } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material";

@Component({
  selector: "app-error-dialog",
  templateUrl: "./error-dialog.component.html",
  styleUrls: ["./error-dialog.component.scss"]
})
export class ErrorDialogComponent implements OnInit {
  constructor(
    private ref: MatDialogRef<ErrorDialogComponent>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      msg: string;
    }
  ) {}

  ngOnInit() {}

  close = () => {
    this.ref.close();
  };
}
