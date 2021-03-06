import * as React from "react";
import { Route, Switch } from "react-router-dom";

import DirectorySection from "./directories/views";
import PageSection from './pages/views'
import UserSection from './users/views'

interface Props {
  match?: any;
}

export const App: React.StatelessComponent<Props> = () => {
  return (
    <Switch>
      <Route path={`/directories`} component={DirectorySection} />
      <Route path={`/pages`} component={PageSection} />
      <Route path={`/users`} component={UserSection} />
    </Switch>
  );
};
export default App;
