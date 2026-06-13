import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { MainTabParamList } from './types';
import { FeedScreen } from '../screens/FeedScreen';
import { ProfileScreen } from '../screens/ProfileScreen';
import { CreateStack } from './CreateStack';
import { TabBar } from './TabBar';

const Tab = createBottomTabNavigator<MainTabParamList>();

export function MainNavigator() {
  return (
    <Tab.Navigator tabBar={(props) => <TabBar {...props} />} screenOptions={{ headerShown: false }}>
      <Tab.Screen name="Feed" component={FeedScreen} />
      <Tab.Screen name="Create" component={CreateStack} />
      <Tab.Screen name="Profile" component={ProfileScreen} />
    </Tab.Navigator>
  );
}
