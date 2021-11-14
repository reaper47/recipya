import { Icon, Menu } from 'semantic-ui-react';
import { Avatar } from './Avatar';

interface HeaderProps {
  isSidebarOpen: boolean;
  handleSidebar(isOpen: boolean): void;
}

export const Header = ({
  isSidebarOpen,
  handleSidebar,
}: HeaderProps): JSX.Element => {
  const makeTitle = (): string => {
    return 'Recipya';
  };

  return (
    <Menu inverted borderless style={{ borderRadius: 0 }}>
      <Menu.Item as='a' header onClick={() => handleSidebar(!isSidebarOpen)}>
        <Icon className={isSidebarOpen ? 'close' : 'bars'} fitted />
      </Menu.Item>

      <Menu.Item>{makeTitle()}</Menu.Item>

      <Menu.Menu position='right'>
        <Menu.Item as='a'>
          <Avatar />
        </Menu.Item>
      </Menu.Menu>
    </Menu>
  );
};
