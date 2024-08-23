using BSC_Main_Backend.Models;

namespace BSC_Main_Backend.Repositories
{
    public interface IAssetRepository
    {
        // Create
        Task<GraphicalAssetModel> AddAssetAsync(GraphicalAssetModel asset);
        Task<IEnumerable<GraphicalAssetModel>> AddAssetsAsync(IEnumerable<GraphicalAssetModel> assets);

        // Read
        Task<GraphicalAssetModel> GetAssetByIdAsync(uint assetId);
        Task<IEnumerable<GraphicalAssetModel>> GetAssetsByIdsAsync(uint[] ids);

        // Update
        Task<GraphicalAssetModel> UpdateAssetAsync(GraphicalAssetModel asset);
        Task<IEnumerable<GraphicalAssetModel>> UpdateAssetsAsync(IEnumerable<GraphicalAssetModel> assets);

        // Delete
        Task<bool> DeleteAssetAsync(uint assetId);
        Task<int> DeleteAssetsAsync(uint[] ids);
    }
}