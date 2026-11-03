import { createSignal, createResource } from 'solid-js';
import { getSettings as getSettingsApi, updateSetting as updateSettingApi } from '../api/settings';

export function useSettings() {
  const [category, setCategory] = createSignal<string | undefined>();

  const [settingsResource, { refetch }] = createResource<Promise<any>, any[]>(() => {
    if (category()) {
      return getSettingsApi(category());
    }
    return getSettingsApi();
  });

  const updateSetting = async (key: string, value: string) => {
    await updateSettingApi(key, { value });
    await refetch();
  };

  const setCategoryFilter = (newCategory: string | undefined) => {
    setCategory(newCategory);
    refetch();
  };

  return {
    settings: settingsResource,
    updateSetting,
    category,
    setCategory: setCategoryFilter,
    isLoading: settingsResource.loading,
  };
}
