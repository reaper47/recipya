import { useEffect } from 'react';
import { Icon, Menu, Sidebar } from 'semantic-ui-react';
import './Drawer.css';

interface DrawerProps {
  isSidebarOpen: boolean;
}

export const Drawer = ({ isSidebarOpen }: DrawerProps): JSX.Element => {
  useEffect(() => {
    console.log(isSidebarOpen);
  }, [isSidebarOpen]);
  const style = {
    paddingLeft: '0',
  };

  return (
    <Sidebar
      as={Menu}
      id='drawer'
      icon='labeled'
      animation='push'
      inverted
      vertical
      visible={isSidebarOpen}
      width='very thin'
      size='mini'
    >
      <Menu.Item as='a' style={style}>
        <Icon name='home' />
        Home
      </Menu.Item>
      <Menu.Item as='a' style={style}>
        <Icon name='gamepad' />
        Games
      </Menu.Item>
      <Menu.Item as='a' style={style}>
        <Icon name='camera' />
        Channels
      </Menu.Item>
    </Sidebar>
  );
};
