import { createSignal, createResource } from 'solid-js';
import { getAccounts, createAccount as createAccountApi, updateAccount as updateAccountApi, deleteAccount as deleteAccountApi } from '../api/account';

export function useAccounts() {
  const [page, setPage] = createSignal(1);
  const [pageSize, setPageSize] = createSignal(10);
  const [searchQuery, setSearchQuery] = createSignal('');
  const [typeFilter, setTypeFilter] = createSignal('');

  const [accountsResource, { refetch }] = createResource(() => {
    const params: { page?: number; per_page?: number; search?: string; type?: string } = {
      page: page(),
      per_page: pageSize(),
    };
    if (searchQuery()) params.search = searchQuery();
    if (typeFilter()) params.type = typeFilter();
    return getAccounts(params);
  });

  const createAccount = async (name: string, parentId?: string | null, accountType?: string) => {
    await createAccountApi({ name, parent_id: parentId ?? null, account_type: accountType });
    refetch();
  };

  const updateAccount = async (id: string, data: { name?: string; parent_id?: string | null; account_type?: string }) => {
    await updateAccountApi(id, data);
    refetch();
  };

  const deleteAccount = async (id: string) => {
    await deleteAccountApi(id);
    refetch();
  };

  const setPageNumber = (newPage: number) => {
    setPage(newPage);
    refetch();
  };

  return {
    accounts: accountsResource,
    refetch,
    createAccount,
    updateAccount,
    deleteAccount,
    searchQuery,
    setSearchQuery,
    typeFilter,
    setTypeFilter,
    page,
    setPage: setPageNumber,
    pageSize,
    setPageSize,
    isLoading: () => accountsResource.loading,
    total: () => accountsResource()?.meta?.total ?? 0,
    totalPages: () => accountsResource()?.meta?.totalPages ?? 0,
  };
}

