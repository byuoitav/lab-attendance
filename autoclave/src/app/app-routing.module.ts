import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";

import { AppComponent } from "./components/app.component";
import { LoginComponent } from "./components/login/login.component";
import { ScreenSaverComponent } from "./components/screen-saver/screen-saver.component";

const routes: Routes = [
  {
    path: "",
    redirectTo: "/login",
    pathMatch: "full"
  },
  {
    path: "",
    component: AppComponent,
    children: [
      {
        path: "login",
        component: LoginComponent
      },
      {
        path: "screensaver",
        component: ScreenSaverComponent
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
