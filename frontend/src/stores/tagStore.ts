import { createEffect, createResource } from "solid-js";
import { createStore, reconcile } from "solid-js/store";
import { DeleteTag, GetTagList, UpdateTag } from "wailsjs/go/main/App";

const [tagList, setTagList] = createStore<{id: number, name: string}[]>([]);

const [dataResource, { refetch }] = createResource(async () => {
    return await GetTagList();
});

createEffect(() => {
    if (dataResource()) setTagList(reconcile(dataResource()!));
})

const submitEditTag = async (
    id: number,
    name: string,
) => {
    try {
        await UpdateTag(id, name);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const deleteTagById = async (id: number) => {
    try {
        await DeleteTag(id);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const getTagById = (id: number) => tagList.find(i => i.id === id);

export {
    tagList,
    submitEditTag,
    deleteTagById,
    getTagById,
};
