import { createEffect, createResource } from "solid-js";
import { createStore, reconcile } from "solid-js/store";
import { CreateBudget, DeleteBudget, GetBudgetList, UpdateBudget } from "wailsjs/go/main/App";
import { domain } from "wailsjs/go/models";

const [budgetList, setBudgetList] = createStore<domain.BudgetDTO[]>([]);

const [dataResource, { refetch }] = createResource(async () => {
    return await GetBudgetList();
});

createEffect(() => {
    if (dataResource()) setBudgetList(reconcile(dataResource()!));
});

const submitNewBudget = async (
    name: string,
    description: string,
    target: string,
    tagId: number,
    frequency: string,
    interval: number,
    dayOfWeek: number,
    dayOfMonth: number,
    startDate: string,
) => {
    const dto: domain.BudgetDTO = {
        id: 0,
        name,
        description,
        target,
        tagId,
        tagName: "",
        frequency,
        interval,
        dayOfWeek,
        dayOfMonth,
        startDate,
        currentSpend: "0",
        nextReset: "",
    };
    await CreateBudget(dto);
    await refetch();
};

const submitEditBudget = async (
    id: number,
    name: string,
    description: string,
    target: string,
    tagId: number,
    frequency: string,
    interval: number,
    dayOfWeek: number,
    dayOfMonth: number,
    startDate: string,
) => {
    try {
        const dto: domain.BudgetDTO = {
            id,
            name,
            description,
            target,
            tagId,
            tagName: "",
            frequency,
            interval,
            dayOfWeek,
            dayOfMonth,
            startDate,
            currentSpend: "0",
            nextReset: "",
        };
        await UpdateBudget(dto);
        await refetch();
    } catch (err) {
        console.log(err);
    }
};

const deleteBudgetById = async (id: number) => {
    try {
        await DeleteBudget(id);
        await refetch();
    } catch (err) {
        console.log(err);
    }
};

const getBudgetById = (id: number) => budgetList.find((i) => i.id === id);

export {
    budgetList,
    submitNewBudget,
    submitEditBudget,
    deleteBudgetById,
    getBudgetById,
    refetch as budgetRefetch,
};
