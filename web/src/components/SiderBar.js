import React, { useContext, useEffect, useMemo, useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { UserContext } from '../context/User';
import { StatusContext } from '../context/Status';
import { useTranslation } from 'react-i18next';

import {
  API,
  getLogo,
  getSystemName,
  isAdmin,
  isMobile,
  showError,
} from '../helpers';
import '../index.css';

import {
  IconCalendarClock, IconChecklistStroked,
  IconComment, IconCommentStroked,
  IconCreditCard,
  IconGift, IconHelpCircle,
  IconHistogram,
  IconHome,
  IconImage,
  IconKey,
  IconLayers,
  IconPriceTag,
  IconSetting,
  IconUser
} from '@douyinfe/semi-icons';
import { Avatar, Dropdown, Layout, Nav, Switch } from '@douyinfe/semi-ui';
import { setStatusData } from '../helpers/data.js';
import { stringToColor } from '../helpers/render.js';
import { useSetTheme, useTheme } from '../context/Theme/index.js';
import { StyleContext } from '../context/Style/index.js';

// HeaderBar Buttons

const SiderBar = () => {
  const { t } = useTranslation();
  const [styleState, styleDispatch] = useContext(StyleContext);
  const [statusState, statusDispatch] = useContext(StatusContext);
  const defaultIsCollapsed =
    localStorage.getItem('default_collapse_sidebar') === 'true';

  const location = useLocation();

  // 根据当前路径设置选中的key
  const getCurrentKey = (pathname) => {
    const path = pathname.split('/')[1] || 'home';
    return path;
  };

  const [selectedKeys, setSelectedKeys] = useState([getCurrentKey(location.pathname)]);
  const [isCollapsed, setIsCollapsed] = useState(defaultIsCollapsed);
  const [chatItems, setChatItems] = useState([]);
  const theme = useTheme();
  const setTheme = useSetTheme();
  
  const controlPaths = ['/channel', '/token', '/user', '/log', '/midjourney', '/setting', '/task', '/playground'];
  const showSidebar = controlPaths.some(path => location.pathname.startsWith(path));

  // 监听路由变化，更新选中状态
  useEffect(() => {
    const currentKey = getCurrentKey(location.pathname);
    setSelectedKeys([currentKey]);
  }, [location.pathname]);

  const routerMap = {
    home: '/',
    channel: '/channel',
    token: '/token',
    user: '/user',
    log: '/log',
    midjourney: '/midjourney',
    setting: '/setting',
    chat: '/chat',
    detail: '/detail',
    pricing: '/pricing',
    task: '/task',
    playground: '/playground',
  };

  const headerButtons = useMemo(
    () => [
      {
        text: t('体验中心'),
        itemKey: 'playground',
        to: '/playground',
        icon: <IconCommentStroked />,
      },
      {
        text: t('渠道'),
        itemKey: 'channel',
        to: '/channel',
        icon: <IconLayers />,
        className: isAdmin() ? 'semi-navigation-item-normal' : 'tableHiddle',
      },
      {
        text: t('令牌'),
        itemKey: 'token',
        to: '/token',
        icon: <IconKey />,
      },
      {
        text: t('用户管理'),
        itemKey: 'user',
        to: '/user',
        icon: <IconUser />,
        className: isAdmin() ? 'semi-navigation-item-normal' : 'tableHiddle',
      },
      {
        text: t('日志'),
        itemKey: 'log',
        to: '/log',
        icon: <IconHistogram />,
      },
      {
        text: t('绘图'),
        itemKey: 'midjourney',
        to: '/midjourney',
        icon: <IconImage />,
        className:
          localStorage.getItem('enable_drawing') === 'true'
            ? 'semi-navigation-item-normal'
            : 'tableHiddle',
      },
      {
        text: t('异步任务'),
        itemKey: 'task',
        to: '/task',
        icon: <IconChecklistStroked />,
        className:
            localStorage.getItem('enable_task') === 'true'
                ? 'semi-navigation-item-normal'
                : 'tableHiddle',
      },
      {
        text: t('设置'),
        itemKey: 'setting',
        to: '/setting',
        icon: <IconSetting />,
      },
    ],
    [
      localStorage.getItem('enable_data_export'),
      localStorage.getItem('enable_drawing'),
      localStorage.getItem('enable_task'),
      t,
    ],
  );

  useEffect(() => {
    let localKey = window.location.pathname.split('/')[1];
    if (localKey === '') {
      localKey = 'home';
    }
    setSelectedKeys([localKey]);
    
    let chatLink = localStorage.getItem('chat_link');
    if (!chatLink) {
        let chats = localStorage.getItem('chats');
        if (chats) {
            // console.log(chats);
            try {
                chats = JSON.parse(chats);
                if (Array.isArray(chats)) {
                    let chatItems = [];
                    for (let i = 0; i < chats.length; i++) {
                        let chat = {};
                        for (let key in chats[i]) {
                            chat.text = key;
                            chat.itemKey = 'chat' + i;
                            chat.to = '/chat/' + i;
                        }
                        // setRouterMap({ ...routerMap, chat: '/chat/' + i })
                        chatItems.push(chat);
                    }
                    setChatItems(chatItems);
                }
            } catch (e) {
                console.error(e);
                showError('聊天数据解析失败')
            }
        }
    }
    
    setIsCollapsed(localStorage.getItem('default_collapse_sidebar') === 'true');
  }, []);

  if (!showSidebar) {
    return null;
  }

  return (
    <>
      <Nav
        style={{ maxWidth: 220, height: '100%' }}
        defaultIsCollapsed={
          localStorage.getItem('default_collapse_sidebar') === 'true'
        }
        isCollapsed={isCollapsed}
        onCollapseChange={(collapsed) => {
          setIsCollapsed(collapsed);
        }}
        selectedKeys={selectedKeys}
        renderWrapper={({ itemElement, isSubNav, isInSubNav, props }) => {
            let chatLink = localStorage.getItem('chat_link');
            if (!chatLink) {
                let chats = localStorage.getItem('chats');
                if (chats) {
                    chats = JSON.parse(chats);
                    if (Array.isArray(chats) && chats.length > 0) {
                        for (let i = 0; i < chats.length; i++) {
                            routerMap['chat' + i] = '/chat/' + i;
                        }
                        if (chats.length > 1) {
                            // delete /chat
                            if (routerMap['chat']) {
                                delete routerMap['chat'];
                            }
                        } else {
                            // rename /chat to /chat/0
                            routerMap['chat'] = '/chat/0';
                        }
                    }
                }
            }
          return (
            <Link
              style={{ textDecoration: 'none' }}
              to={routerMap[props.itemKey]}
            >
              {itemElement}
            </Link>
          );
        }}
        items={headerButtons}
        onSelect={(key) => {
          if (key.itemKey.toString().startsWith('chat')) {
            styleDispatch({ type: 'SET_INNER_PADDING', payload: false });
          } else {
            styleDispatch({ type: 'SET_INNER_PADDING', payload: true });
          }
          setSelectedKeys([key.itemKey]);
        }}
        footer={
          <>
          </>
        }
      >
        <Nav.Footer collapseButton={true}></Nav.Footer>
      </Nav>
    </>
  );
};

export default SiderBar;
