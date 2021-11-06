import { Menu } from 'semantic-ui-react';
import { Avatar } from './Avatar';
import { Drawer } from './Drawer';

export const Header = (): JSX.Element => {
  const makeTitle = (): string => {
    return 'Recipya';
  };

  return (
    <Menu inverted borderless style={{ borderRadius: 0 }}>
      <Menu.Item as='a' header>
        <Drawer />
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
