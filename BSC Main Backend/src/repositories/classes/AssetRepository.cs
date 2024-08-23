using BSC_Main_Backend.Models;
using Microsoft.EntityFrameworkCore;

namespace BSC_Main_Backend.Repositories
{
    public class AssetRepository : IAssetRepository
    {
        private readonly DbContext _context;

        public AssetRepository(DbContext context)
        {
            _context = context;
        }

        // Create
        public async Task<GraphicalAssetModel> AddAssetAsync(GraphicalAssetModel asset)
        {
            _context.GraphicalAssets.Add(asset);
            await _context.SaveChangesAsync();
            return asset;
        }

        public async Task<IEnumerable<GraphicalAssetModel>> AddAssetsAsync(IEnumerable<GraphicalAssetModel> assets)
        {
            _context.GraphicalAssets.AddRange(assets);
            await _context.SaveChangesAsync();
            return assets;
        }

        // Read
        public async Task<GraphicalAssetModel> GetAssetByIdAsync(uint assetId)
        {
            return await _context.GraphicalAssets.FindAsync(assetId);
        }

        public async Task<IEnumerable<GraphicalAssetModel>> GetAssetsByIdsAsync(uint[] ids)
        {
            return await _context.GraphicalAssets
                .Where(asset => ids.Contains((uint)asset.Id))
                .ToListAsync();
        }

        // Update
        public async Task<GraphicalAssetModel> UpdateAssetAsync(GraphicalAssetModel asset)
        {
            var existingAsset = await _context.GraphicalAssets.FindAsync(asset.Id);
            if (existingAsset == null)
            {
                return null;
            }

            _context.Entry(existingAsset).CurrentValues.SetValues(asset);
            await _context.SaveChangesAsync();
            return existingAsset;
        }

        public async Task<IEnumerable<GraphicalAssetModel>> UpdateAssetsAsync(IEnumerable<GraphicalAssetModel> assets)
        {
            var assetIds = assets.Select(a => a.Id).ToArray();
            var existingAssets = await _context.GraphicalAssets.Where(a => assetIds.Contains(a.Id)).ToListAsync();

            if (existingAssets == null || existingAssets.Count != assets.Count())
            {
                return null; // Handle cases where some assets were not found
            }

            foreach (var asset in assets)
            {
                var existingAsset = existingAssets.FirstOrDefault(a => a.Id == asset.Id);
                if (existingAsset != null)
                {
                    _context.Entry(existingAsset).CurrentValues.SetValues(asset);
                }
            }

            await _context.SaveChangesAsync();
            return existingAssets;
        }

        // Delete
        public async Task<bool> DeleteAssetAsync(uint assetId)
        {
            var asset = await _context.GraphicalAssets.FindAsync(assetId);
            if (asset == null)
            {
                return false;
            }

            _context.GraphicalAssets.Remove(asset);
            await _context.SaveChangesAsync();
            return true;
        }

        public async Task<int> DeleteAssetsAsync(uint[] ids)
        {
            var assets = await _context.GraphicalAssets.Where(asset => ids.Contains((uint)asset.Id)).ToListAsync();
            if (assets == null || !assets.Any())
            {
                return 0; // No assets found to delete
            }

            _context.GraphicalAssets.RemoveRange(assets);
            return await _context.SaveChangesAsync();
        }
    }
}
