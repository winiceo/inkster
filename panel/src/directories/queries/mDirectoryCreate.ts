import gql from "graphql-tag";

export interface DirectoryCreateVariables {
  name: string;
  parentId?: string;
}
const mDirectoryCreate = gql`
  mutation DirectoryCreate($name: String!, $parentId: ID) {
    createDirectory(input: { name: $name, parentId: $parentId }) {
      errors {
        field
        code
      }
      directory {
        id
        createdAt
        updatedAt
        name
        isPublished
        parent {
          id
        }
        pages {
          id
        }
      }
    }
  }
`;
export default mDirectoryCreate;
