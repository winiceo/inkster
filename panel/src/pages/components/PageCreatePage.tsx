import * as React from "react";
import { Plus } from "react-feather";
import withStyles from "react-jss";

import PageHeader from "../../components/PageHeader";
import Container from "../../components/Container";
import IconButton from "../../components/IconButton";
import ActionDialog from "../../components/ActionDialog";
import Input from "../../components/Input";
import Form from "../../components/Form";
import FormSave from "../../components/FormSave";
import Toggle from "../../components/Toggle";
import { ViewProps, FormViewProps } from "../../";
import i18n from "../../i18n";
import PageProperties from "./PageProperties";
import PageFieldProperties from "./PageFieldProperties";

interface PageField {
  id: string;
  name: string;
  type: string;
  value: string;
}
export interface FormData {
  name: string;
  slug: string;
  addFields: PageField[];
  removeFields: string[];
}
interface Props extends ViewProps, FormViewProps<FormData> {
  onUpload: (
    cb: (event: React.ChangeEvent<any>) => void
  ) => (event: React.ChangeEvent<any>) => void;
}

const decorate = withStyles(
  (theme: any) => ({
    root: {
      display: "grid" as "grid",
      gridColumnGap: theme.spacing + "px",
      gridTemplateColumns: "2fr 1fr"
    }
  }),
  { displayName: "PageCreatePage" }
);
export const PageCreatePage = decorate<Props>(
  ({
    classes,
    disabled,
    loading,
    title,
    transaction,
    onBack,
    onUpload,
    onSubmit
  }) => {
    const initialForm = {
      name: "",
      slug: "",
      addFields: [] as PageField[],
      removeFields: [] as string[]
    };
    return (
      <Form
        initial={initialForm}
        onSubmit={onSubmit}
      >
        {({ change, data, hasChanged, submit }) => {
          const handleFieldAdd = (field: { type: string }) => {
            change({
              target: {
                name: "addFields",
                value: [
                  ...data.addFields,
                  {
                    type: field.type,
                    id: "new-" + data.addFields.length,
                    name: "",
                    value: ""
                  }
                ]
              }
            } as any);
          };
          const handleFieldRemove = (name: string, id: string) => () => {
            change({
              target: {
                name,
                value: data[name].filter((f: PageField) => f.id !== id)
              }
            } as any);
            if (name === "fields") {
              change({
                target: {
                  name: "removeFields",
                  value: [id, ...data.removeFields]
                }
              } as any);
            }
          };
          const handleChange = (name: string, id: string) => (
            event: React.ChangeEvent<any>
          ) =>
            change({
              target: {
                name,
                value: data[name].map(
                  (f: PageField) =>
                    f.id === id
                      ? { ...f, [event.target.name]: event.target.value }
                      : f
                )
              }
            } as any);
          return (
            <Toggle>
              {(openedRemoveDialog, { toggle: toggleRemoveDialog }) => (
                <Toggle>
                  {(openedFieldAddDialog, { toggle: toggleFieldAddDialog }) => (
                    <Toggle>
                      {(
                        openedFieldRemoveDialog,
                        { toggle: toggleFieldRemoveDialog }
                      ) => (
                        <>
                          <Container width="md">
                            <PageHeader onBack={onBack} title={title}>
                              <IconButton
                                disabled={disabled || loading}
                                icon={Plus}
                                onClick={toggleFieldAddDialog}
                              />
                            </PageHeader>
                            <div className={classes.root}>
                              <div>
                                <PageProperties
                                  data={data}
                                  disabled={disabled || loading}
                                  onChange={change}
                                />
                                {data.addFields.map((field, index) => (
                                  <PageFieldProperties
                                    data={field}
                                    key={field.id + index}
                                    name="addFields"
                                    onChange={handleChange(
                                      "addFields",
                                      field.id
                                    )}
                                    onDelete={handleFieldRemove(
                                      "addFields",
                                      field.id
                                    )}
                                    onUpload={onUpload}
                                  />
                                ))}
                              </div>
                            </div>
                            <FormSave
                              disabled={disabled || loading || !hasChanged}
                              onConfirm={submit}
                              variant={transaction}
                            />
                          </Container>
                          {!disabled &&
                            !loading && (
                              <>
                                <Form
                                  initial={{ type: "text" }}
                                  onSubmit={handleFieldAdd}
                                >
                                  {({
                                    change: handleAddFieldChange,
                                    data: addFieldData,
                                    submit: addField
                                  }) => (
                                    <ActionDialog
                                      show={openedFieldAddDialog}
                                      size="xs"
                                      title={i18n.t("Add page field")}
                                      onClose={toggleFieldAddDialog}
                                      onConfirm={addField as () => void}
                                    >
                                      <Input
                                        name="type"
                                        label={i18n.t("Type")}
                                        value={addFieldData.type}
                                        onChange={handleAddFieldChange}
                                        type="select"
                                      >
                                        <>
                                          <option value="text">
                                            {i18n.t("Short text")}
                                          </option>
                                          <option value="longText">
                                            {i18n.t("Long text")}
                                          </option>
                                          <option value="image">
                                            {i18n.t("Image")}
                                          </option>
                                          <option value="file">
                                            {i18n.t("File")}
                                          </option>
                                        </>
                                      </Input>
                                    </ActionDialog>
                                  )}
                                </Form>
                              </>
                            )}
                        </>
                      )}
                    </Toggle>
                  )}
                </Toggle>
              )}
            </Toggle>
          );
        }}
      </Form>
    );
  }
);
export default PageCreatePage;
