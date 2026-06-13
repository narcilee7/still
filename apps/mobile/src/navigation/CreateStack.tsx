import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { CreateStackParamList } from './types';
import { CreateSelectScreen } from '../screens/create/CreateSelectScreen';
import { CreateEditScreen } from '../screens/create/CreateEditScreen';
import { CreateSuccessScreen } from '../screens/create/CreateSuccessScreen';

const Stack = createNativeStackNavigator<CreateStackParamList>();

export function CreateStack() {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false }}>
      <Stack.Screen name="CreateSelect" component={CreateSelectScreen} />
      <Stack.Screen name="CreateEdit" component={CreateEditScreen} />
      <Stack.Screen name="CreateSuccess" component={CreateSuccessScreen} />
    </Stack.Navigator>
  );
}
