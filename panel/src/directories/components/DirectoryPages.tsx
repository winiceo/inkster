import * as React from "react";
import { Panel } from "react-bootstrap";
import { FileText, Plus } from "react-feather";

import { PaginatedListProps } from "../..";
import ListElement from "../../components/ListElement";
import IconButton from "../../components/IconButton";
import Paginator from "../../components/Paginator";
import i18n from "../../i18n";

interface Props extends PaginatedListProps {
  disabled: boolean;
  pages?: Array<{
    id: string;
    name?: string;
  }>;
  onAdd: () => void;
}

export const DirectoryRootList: React.StatelessComponent<Props> = ({
  disabled,
  pages,
  pageInfo,
  onAdd,
  onNextPage,
  onPreviousPage,
  onRowClick
}) => (
  <Panel>
    <Panel.Heading>
      <Panel.Title>{i18n.t("Pages")}</Panel.Title>
  <IconButton disabled={disabled} icon={Plus} onClick={onAdd} />
    </Panel.Heading>
    <Panel.Body>
      {pages ? (
        pages.length > 0 ? (
          pages.map(page => (
            <ListElement
              disabled={disabled}
              title={page.name}
              onClick={onRowClick(page.id)}
              icon={FileText}
            />
          ))
        ) : (
          i18n.t("No pages found")
        )
      ) : (
        <ListElement disabled={disabled} icon={FileText} />
      )}
    </Panel.Body>
    <Panel.Footer>
      <Paginator
        pageInfo={pageInfo}
        onNextPage={onNextPage}
        onPreviousPage={onPreviousPage}
      />
    </Panel.Footer>
  </Panel>
);
export default DirectoryRootList;
