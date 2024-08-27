using BSC_Main_Backend.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace BSC_Main_Backend.Data
{
    public interface IGraphicalAssetRepository
    {
        Task<IEnumerable<GraphicalAsset>> GetAllGraphicalAssetsAsync();
        Task<GraphicalAsset> GetGraphicalAssetByIdAsync(int id);
        Task CreateGraphicalAssetAsync(GraphicalAsset graphicalAsset);
        Task UpdateGraphicalAssetAsync(GraphicalAsset graphicalAsset);
        Task DeleteGraphicalAssetAsync(int id);
    }
}