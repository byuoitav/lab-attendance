import { Component, OnInit, Inject } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material";

@Component({
  selector: "app-error-dialog",
  templateUrl: "./error-dialog.component.html",
  styleUrls: ["./error-dialog.component.scss"]
})
export class ErrorDialogComponent implements OnInit {
  constructor(
    public ref: MatDialogRef<ErrorDialogComponent>,
    @Inject(MAT_DIALOG_DATA)
    public data: {
      msg: string;
    }
  ) {
    setTimeout(() => {
      this.ref.close();
    }, 3000);
  }

  ngOnInit() {}

  close = () => {
    this.ref.close();
  };
}
