import * as React from "react";
import { Mutation } from "react-apollo";
import { ApolloError } from "apollo-client";

import { TransactionState } from "../../";
import urls from "../../urls";
import Navigator from "../../components/Navigator";
import Notificator, { NotificationType } from "../../components/Notificator";
import PageCreatePage, { FormData } from "../components/PageCreatePage";
import mPageCreate, {
  result as PageCreateResult,
  variables as PageCreateVariables,
} from "../queries/mPageCreate";
import i18n from "../../i18n";
import { WithUpload } from "../../UploadProvider";

interface Props {
  directory: string;
}
interface State {
  transaction: TransactionState;
}
export class PageCreate extends React.Component<Props, State> {
  state = {
    transaction: "default" as "default",
  };

  handleSubmitSuccess = (cb: () => void) => {
    this.setState({ transaction: "success" });
    setTimeout(cb, 3000);
  };
  handleSubmitError = (event: ApolloError) => {
    this.setState({ transaction: "error" });
    setTimeout(() => this.setState({ transaction: "default" }), 3000);
  };

  render() {
    const { directory } = this.props;
    return (
      <Notificator>
        {notify => (
          <Navigator>
            {navigate => {
              const handleBack = () =>
                navigate(urls.directoryDetails(directory));
              const handleError = () =>
                notify({
                  text: i18n.t("Something has gone wrong", {
                    context: "notification",
                  }),
                  type: NotificationType.ERROR,
                });
              const handleSubmitSuccess = (data: {
                createPage: PageCreateResult;
              }) => {
                if (
                  data.createPage &&
                  data.createPage.errors &&
                  data.createPage.errors.length > 0
                ) {
                  handleError();
                }
                notify({
                  text: i18n.t("Page created", {
                    context: "notification",
                  }),
                });
                navigate(urls.pageDetails(data.createPage.page.id));
              };
              return (
                <WithUpload>
                  {uploadFile => {
                    const handleUpload = (onChange: any) => (
                      event: React.ChangeEvent<any>,
                    ) => {
                      uploadFile.uploadFile({
                        file: event.target.files[0],
                        onSuccess: filename =>
                          onChange({
                            target: {
                              name: "value",
                              value: filename,
                            },
                          } as any),
                        onError: () => console.log("not ok"),
                      });
                    };
                    return (
                      <Mutation
                        mutation={mPageCreate}
                        onCompleted={handleSubmitSuccess}
                      >
                        {(createPage, { data, loading }) => {
                          const handleSubmit = (formData: FormData) =>
                            createPage({
                              variables: {
                                name: formData.name,
                                parentId: directory,
                                fields: formData.addFields,
                              } as PageCreateVariables,
                            });
                          return (
                            <PageCreatePage
                              disabled={loading}
                              loading={loading}
                              title={i18n.t("Create new page")}
                              transaction={this.state.transaction}
                              onBack={handleBack}
                              onUpload={handleUpload}
                              onSubmit={handleSubmit}
                            />
                          );
                        }}
                      </Mutation>
                    );
                  }}
                </WithUpload>
              );
            }}
          </Navigator>
        )}
      </Notificator>
    );
  }
}
export default PageCreate;
